package common

import "errors"

func GetMsg(key int) string {
	return message[key]
}

func ReturnErr(code int) error {
	return errors.New(message[code])
}

func Result(code int, data any, msg ...any) map[string]any {
	if data == nil || data == "" {
		data = map[string]any{}
	}
	res := map[string]any{
		"code": code,
		"data": data,
	}
	if len(msg) > 0 {
		// if err, ok := msg[0].(error); ok {
		// 	res["message"] = err.Error()
		// } else {
		// 	res["message"] = msg[0]
		// }
		message := ""
		for _, v := range msg {
			if err, ok := v.(error); ok {
				message += err.Error()
			} else if str, ok := v.(string); ok {
				message += str
			}
		}
		res["message"] = message

	} else if v, ok := message[code]; ok {
		res["message"] = v
	}
	return res
}
