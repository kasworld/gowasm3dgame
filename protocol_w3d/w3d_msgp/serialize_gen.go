// Code generated by "genprotocol -ver=de39b3963afaad4ed2557809c9208a5676bf9321c0823e0dac15bc4db6d51552 -basedir=. -prefix=w3d -statstype=int"

package w3d_msgp

import (
	"fmt"

	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_idcmd"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_idnoti"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_packet"
	"github.com/tinylib/msgp/msgp"
)

// MarshalBodyFn marshal body and append to oldBufferToAppend
// and return newbuffer, body type, error
func MarshalBodyFn(body interface{}, oldBuffToAppend []byte) ([]byte, byte, error) {
	newBuffer, err := body.(msgp.Marshaler).MarshalMsg(oldBuffToAppend)
	return newBuffer, 0, err
}

func UnmarshalPacket(h w3d_packet.Header, bodyData []byte) (interface{}, error) {
	switch h.FlowType {
	case w3d_packet.Request:
		if int(h.Cmd) >= len(ReqUnmarshalMap) {
			return nil, fmt.Errorf("unknown request command: %v %v",
				h.FlowType, w3d_idcmd.CommandID(h.Cmd))
		}
		return ReqUnmarshalMap[h.Cmd](h, bodyData)

	case w3d_packet.Response:
		if int(h.Cmd) >= len(RspUnmarshalMap) {
			return nil, fmt.Errorf("unknown response command: %v %v",
				h.FlowType, w3d_idcmd.CommandID(h.Cmd))
		}
		return RspUnmarshalMap[h.Cmd](h, bodyData)

	case w3d_packet.Notification:
		if int(h.Cmd) >= len(NotiUnmarshalMap) {
			return nil, fmt.Errorf("unknown notification command: %v %v",
				h.FlowType, w3d_idcmd.CommandID(h.Cmd))
		}
		return NotiUnmarshalMap[h.Cmd](h, bodyData)
	}
	return nil, fmt.Errorf("unknown packet FlowType %v", h.FlowType)
}

var ReqUnmarshalMap = [...]func(h w3d_packet.Header, bodyData []byte) (interface{}, error){
	w3d_idcmd.Invalid:   unmarshal_ReqInvalid,
	w3d_idcmd.Login:     unmarshal_ReqLogin,
	w3d_idcmd.Heartbeat: unmarshal_ReqHeartbeat,
	w3d_idcmd.Chat:      unmarshal_ReqChat,
}

var RspUnmarshalMap = [...]func(h w3d_packet.Header, bodyData []byte) (interface{}, error){
	w3d_idcmd.Invalid:   unmarshal_RspInvalid,
	w3d_idcmd.Login:     unmarshal_RspLogin,
	w3d_idcmd.Heartbeat: unmarshal_RspHeartbeat,
	w3d_idcmd.Chat:      unmarshal_RspChat,
}

var NotiUnmarshalMap = [...]func(h w3d_packet.Header, bodyData []byte) (interface{}, error){
	w3d_idnoti.Invalid:   unmarshal_NotiInvalid,
	w3d_idnoti.StageInfo: unmarshal_NotiStageInfo,
	w3d_idnoti.StatsInfo: unmarshal_NotiStatsInfo,
	w3d_idnoti.StageChat: unmarshal_NotiStageChat,
}

func unmarshal_ReqInvalid(h w3d_packet.Header, bodyData []byte) (interface{}, error) {
	var args w3d_obj.ReqInvalid_data
	if _, err := args.UnmarshalMsg(bodyData); err != nil {
		return nil, err
	}
	return &args, nil
}

func unmarshal_RspInvalid(h w3d_packet.Header, bodyData []byte) (interface{}, error) {
	var args w3d_obj.RspInvalid_data
	if _, err := args.UnmarshalMsg(bodyData); err != nil {
		return nil, err
	}
	return &args, nil
}

func unmarshal_ReqLogin(h w3d_packet.Header, bodyData []byte) (interface{}, error) {
	var args w3d_obj.ReqLogin_data
	if _, err := args.UnmarshalMsg(bodyData); err != nil {
		return nil, err
	}
	return &args, nil
}

func unmarshal_RspLogin(h w3d_packet.Header, bodyData []byte) (interface{}, error) {
	var args w3d_obj.RspLogin_data
	if _, err := args.UnmarshalMsg(bodyData); err != nil {
		return nil, err
	}
	return &args, nil
}

func unmarshal_ReqHeartbeat(h w3d_packet.Header, bodyData []byte) (interface{}, error) {
	var args w3d_obj.ReqHeartbeat_data
	if _, err := args.UnmarshalMsg(bodyData); err != nil {
		return nil, err
	}
	return &args, nil
}

func unmarshal_RspHeartbeat(h w3d_packet.Header, bodyData []byte) (interface{}, error) {
	var args w3d_obj.RspHeartbeat_data
	if _, err := args.UnmarshalMsg(bodyData); err != nil {
		return nil, err
	}
	return &args, nil
}

func unmarshal_ReqChat(h w3d_packet.Header, bodyData []byte) (interface{}, error) {
	var args w3d_obj.ReqChat_data
	if _, err := args.UnmarshalMsg(bodyData); err != nil {
		return nil, err
	}
	return &args, nil
}

func unmarshal_RspChat(h w3d_packet.Header, bodyData []byte) (interface{}, error) {
	var args w3d_obj.RspChat_data
	if _, err := args.UnmarshalMsg(bodyData); err != nil {
		return nil, err
	}
	return &args, nil
}

func unmarshal_NotiInvalid(h w3d_packet.Header, bodyData []byte) (interface{}, error) {
	var args w3d_obj.NotiInvalid_data
	if _, err := args.UnmarshalMsg(bodyData); err != nil {
		return nil, err
	}
	return &args, nil
}

func unmarshal_NotiStageInfo(h w3d_packet.Header, bodyData []byte) (interface{}, error) {
	var args w3d_obj.NotiStageInfo_data
	if _, err := args.UnmarshalMsg(bodyData); err != nil {
		return nil, err
	}
	return &args, nil
}

func unmarshal_NotiStatsInfo(h w3d_packet.Header, bodyData []byte) (interface{}, error) {
	var args w3d_obj.NotiStatsInfo_data
	if _, err := args.UnmarshalMsg(bodyData); err != nil {
		return nil, err
	}
	return &args, nil
}

func unmarshal_NotiStageChat(h w3d_packet.Header, bodyData []byte) (interface{}, error) {
	var args w3d_obj.NotiStageChat_data
	if _, err := args.UnmarshalMsg(bodyData); err != nil {
		return nil, err
	}
	return &args, nil
}
