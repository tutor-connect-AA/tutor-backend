package api_ports

import (
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
)

/*A client should have the following services:
1. Register
2. Login
3. Be listed [A list of clients](paginated)
4. Post a job //Create job service
5. Update profile
6. Send hiring request //hiring request service
7. Profile should be shown in detail based on id
*/

type APIPort interface {
	RegisterClient(usr domain.Client) (*domain.Client, error)
	GetClientById(id string) (*domain.Client, error)
	GetListOfClients() ([]*domain.Client, error)
	UpdateClientProfile(updatedClt domain.Client) error //Takes in the updated client and returns the id of the client if success
	LoginClient(username, password string) (bool, error)
}
