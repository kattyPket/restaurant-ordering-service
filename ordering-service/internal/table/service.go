package table

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	repo       *Repository
	sessionTTL time.Duration
}

var (
	ErrTableNotFound         = errors.New("table not found")
	ErrTableNotUsable        = errors.New("table is under maintenance")
	ErrSessionExpired        = errors.New("session expired or invalid")
	ErrTableAwaitingCheckout = errors.New("table is awaiting checkout")
	ErrSessionTableMismatch  = errors.New("session does not belong to table")
)

func NewService(repo *Repository, sessionTTL time.Duration) *Service {
	return &Service{
		repo:       repo,
		sessionTTL: sessionTTL,
	}
}

func (s *Service) StartSession(qrToken string) (*TableSession, *DiningTable, error) {
	var finalSession *TableSession
	var finalTable *DiningTable

	// ເລີ່ມ Transaction
	err := s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		// 1. ຫາໂຕະ ແລະ ລັອກແຖວໄວ້ (FOR UPDATE)
		dt, err := s.repo.FindTableByQRTokenWithLock(tx, qrToken)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrTableNotFound
			}
			return err
		}

		if dt.Status == "maintenance" {
			return ErrTableNotUsable
		}
		finalTable = dt

		// 2. ກວດຫາ Session ທີ່ເຄີຍເປີດໄວ້
		existingSession, err := s.repo.FindActiveSessionByTableID(tx, dt.ID.String())

		// ແຍກ ErrRecordNotFound ອອກຈາກ Database Error ຢ່າງຊັດເຈນ
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err // DB error ອື່ນໆ
		}

		if existingSession != nil {
			// ມີ Session ເກົ່າ
			if dt.Status == "occupied" && time.Now().After(existingSession.ExpiresAt) {
				return ErrTableAwaitingCheckout
			}
			// An active session always wins: QR retries receive the same session.
			finalSession = existingSession
			return nil
		}
		if dt.Status != "available" {
			return ErrTableAwaitingCheckout
		}

		// 3. ສ້າງ Session ໃໝ່
		tok, err := generateToken()
		if err != nil {
			return err
		}

		newSession := &TableSession{
			TableID:   dt.ID,
			Token:     tok,
			CreateAt:  time.Now(),
			ExpiresAt: time.Now().Add(s.sessionTTL),
		}

		if err := s.repo.CreateSession(tx, newSession); err != nil {
			return err
		}
		finalSession = newSession

		// 4. ອັບເດດສະຖານະໂຕະ ຖ້າຍັງວ່າງຢູ່ (ຫ້າມ Ignore Error ເດັດຂາດ)
		if dt.Status == "available" {
			if err := s.repo.UpdateTableStatus(tx, dt.ID.String(), "occupied"); err != nil {
				return err
			}
		}

		return nil // ບັນທຶກສຳເລັດ (Commit)
	})

	if err != nil {
		return nil, nil, err
	}

	return finalSession, finalTable, nil
}

func (s *Service) CloseSessionForCheckout(sessionID, tableID uuid.UUID) error {
	return s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		dt, err := s.repo.FindTableByIDWithLock(tx, tableID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrTableNotFound
			}
			return err
		}
		sess, err := s.repo.FindSessionByID(tx, sessionID)
		if err != nil {
			return err
		}
		if sess.TableID != dt.ID {
			return ErrSessionTableMismatch
		}
		if sess.ClosedAt != nil {
			return nil
		}
		now := time.Now()
		if err := tx.Model(&TableSession{}).Where("id = ?", sessionID).Update("closed_at", now).Error; err != nil {
			return err
		}
		return s.repo.UpdateTableStatus(tx, dt.ID.String(), "available")
	})
}

func (s *Service) ValidationToken(tokenStr string) (*TableSession, error) {
	sess, err := s.repo.FindSessionByToken(tokenStr)
	if err != nil {
		return nil, ErrSessionExpired
	}

	if sess.ClosedAt != nil || time.Now().After(sess.ExpiresAt) {
		return nil, ErrSessionExpired
	}
	return sess, nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
