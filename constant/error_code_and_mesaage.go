package constant

// 错误码规则
// code = methodID + message code
// 公共错误的methodID 为 0
// 例如：method 00 公共错误，请求格式不正确。  00 + 001 = 00001
// 例如：method 11 产生的错误。     11 + 001 = 11001

const SUCCESS = "00000"
const ERR_JSON_MARSHALLER_FAIL = "00001"
const ERR_GET_TCP_USER_FAIL = "00003"
const ERR_TOKEN = "10001"
const GAME_ERROR = "20001"
