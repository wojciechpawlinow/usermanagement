package domain

import (
	"fmt"

	"github.com/google/uuid"
)

// ID is a custom type that abstracts the underlying identifier implementation
type ID struct {
	value uuid.UUID
}

func NewID() ID {
	return ID{value: uuid.New()}
}

func ParseID(id string) (ID, error) {
	uuidValue, err := uuid.Parse(id)
	if err != nil {
		return ID{}, err
	}

	return ID{value: uuidValue}, nil
}

func (id ID) String() string {
	return id.value.String()
}

func (id ID) Equal(other ID) bool {
	return id.value == other.value
}

func (id ID) IsEmpty() bool {
	return id.value == uuid.Nil
}

func (id ID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, id.value.String())), nil
}

func (id ID) UnmarshalJSON(data []byte) error {
	parsedID, err := uuid.Parse(string(data[1 : len(data)-1]))
	if err != nil {
		return err
	}

	id.value = parsedID

	return nil
}
