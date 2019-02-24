package notify

import (
	"errors"
	"github.com/checkr/flagr/pkg/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestIntegration implements a Notifier for testing purposes.
type TestIntegration struct {
	fakeErr error
	callCount int
}

// Notify handles notifications for a TestIntegration
func (n *TestIntegration) Notify(f *entity.Flag, b notify, i itemType) error {
	n.callCount++
	return n.fakeErr
}

func TestNotifyAll(t *testing.T) {
	t.Run("we return early when the flagID is not found", func(t *testing.T) {
		db := entity.NewTestDB()
		defer db.Close()
		All(db, 1, TOGGLED, FLAG)
	})

	t.Run("nothing bad happens if we have the flag, but no configured integrations", func(t *testing.T) {
		f := entity.GenFixtureFlag()
		db := entity.PopulateTestDB(f)
		defer db.Close()
		All(db, f.ID, TOGGLED, FLAG)
	})

	t.Run("nothing bad happens if an integration fails to deliver and returns an error", func(t *testing.T) {
		f := entity.GenFixtureFlag()
		db := entity.PopulateTestDB(f)
		defer db.Close()

		notifier := &TestIntegration{fakeErr: errors.New("failed to notify testcase")}
		integration := Integration{notifier: notifier, name: "test"}
		Integrations = append(Integrations, integration)
		All(db, f.ID, TOGGLED, FLAG)
		assert.Equal(t, 1, notifier.callCount)
	})
}