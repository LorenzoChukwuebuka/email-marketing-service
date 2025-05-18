package payloads 

type UserDetailsPayload struct {
    Details map[string]interface{} `json:"details"`
}