package policies

import (
	"context"
	"time"

	"github.com/mainflux/mainflux/internal/apiutil"
)

// PolicyTypes contains a list of the available policy types currently supported
// They are arranged in the following order based on their priority:
//
// Group policies
//  1. g_add - adds a member to a group
//  2. g_delete - delete a group
//  3. g_update - update a group
//  4. g_list - list groups and their members
//
// Client policies
//  5. c_delete - delete a client
//  6. c_update - update a client
//  8. c_list - list clients
//
// Message policies
//  9. m_write - write a message
//  10. m_read - read a message
var PolicyTypes = []string{"g_add", "g_delete", "g_update", "g_list", "c_delete", "c_update", "c_list", "m_write", "m_read"}

// Policy represents an argument struct for making a policy related function calls.
type Policy struct {
	OwnerID   string    `json:"owner_id"`
	Subject   string    `json:"subject"`
	Object    string    `json:"object"`
	Actions   []string  `json:"actions"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	UpdatedBy string    `json:"updated_by,omitempty"`
}

// PolicyPage contains a page of policies.
type PolicyPage struct {
	Page
	Policies []Policy
}

// PolicyRepository specifies an account persistence API.
type PolicyRepository interface {
	// Save creates a policy for the given Subject, so that, after
	// Save, `Subject` has a `relation` on `group_id`. Returns a non-nil
	// error in case of failures.
	Save(ctx context.Context, p Policy) error

	// CheckAdmin checks if the user is an admin user
	CheckAdmin(ctx context.Context, id string) error

	// Evaluate is used to evaluate if you have the correct permissions.
	// We evaluate if we are in the same group first then evaluate if the
	// object has that action over the subject
	Evaluate(ctx context.Context, entityType string, p Policy) error

	// Update updates the policy type.
	Update(ctx context.Context, p Policy) error

	// Retrieve retrieves policy for a given input.
	Retrieve(ctx context.Context, pm Page) (PolicyPage, error)

	// Delete deletes the policy
	Delete(ctx context.Context, p Policy) error
}

// PolicyService represents a authorization service. It exposes
// functionalities through `auth` to perform authorization.
type PolicyService interface {
	// Authorize checks authorization of the given `subject`. Basically,
	// Authorize verifies that Is `subject` allowed to `relation` on
	// `object`. Authorize returns a non-nil error if the subject has
	// no relation on the object (which simply means the operation is
	// denied).
	Authorize(ctx context.Context, entityType string, p Policy) error

	// AddPolicy creates a policy for the given subject, so that, after
	// AddPolicy, `subject` has a `relation` on `object`. Returns a non-nil
	// error in case of failures.
	AddPolicy(ctx context.Context, token string, p Policy) error

	// UpdatePolicy updates policies based on the given policy structure.
	UpdatePolicy(ctx context.Context, token string, p Policy) error

	// ListPolicy lists policies based on the given policy structure.
	ListPolicy(ctx context.Context, token string, pm Page) (PolicyPage, error)

	// DeletePolicy removes a policy.
	DeletePolicy(ctx context.Context, token string, p Policy) error
}

// Validate returns an error if policy representation is invalid.
func (p Policy) Validate() error {
	if p.Subject == "" {
		return apiutil.ErrMissingPolicySub
	}
	if p.Object == "" {
		return apiutil.ErrMissingPolicyObj
	}
	if len(p.Actions) == 0 {
		return apiutil.ErrMalformedPolicyAct
	}
	for _, p := range p.Actions {
		if ok := ValidateAction(p); !ok {
			return apiutil.ErrMalformedPolicyAct
		}
	}
	return nil
}

// ValidateAction check if the action is in policies
func ValidateAction(act string) bool {
	for _, v := range PolicyTypes {
		if v == act {
			return true
		}
	}
	return false

}