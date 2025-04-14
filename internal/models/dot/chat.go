package dot

type Address struct {
	Name     string `json:"name"`
	DeviceId int    `json:"deviceId"`
}

type SignedPreKey struct {
	Id        int    `json:"id"`
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
}

type PreKey struct {
	Id        int    `json:"id"`
	PublicKey string `json:"publicKey"`
}

type SignalData struct {
	Address        Address      `json:"address"`
	RegistrationId int          `json:"registrationId"`
	IdentityKey    string       `json:"identityKey"`
	SignedPreKey   SignedPreKey `json:"signedPreKey"`
	PreKey         PreKey       `json:"preKey"`
}

type JoinRoomData struct {
	RoomUUID string `json:"room_uuid" binding:"required"`
	UserUUID string `json:"user_uuid" binding:"required"`
	Password string `json:"password"`
}

type CreateRoomData struct {
	UserUUID string `json:"user_uuid" binding:"required"`
	RoomName string `json:"room_name" binding:"required"`
	Password string `json:"password"`
}
