package server

import (
	"encoding/json"
	"io"
	"net"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type StreamedResponse struct {
	MessageType string `json:"Type"`
	MessageData string `json:"Data"`
	PercentDone uint8  `json:"Percent"`
}

func WriteJSONInfo(conn net.Conn, data string, percent uint8) error {
	return WriteJSON(conn, &StreamedResponse{
		MessageType: "info",
		MessageData: data,
		PercentDone: percent,
	})
}

func WriteJSONSuccess(conn net.Conn, data string) error {
	return WriteJSON(conn, &StreamedResponse{
		MessageType: "success",
		MessageData: data,
		PercentDone: 100,
	})
}

func WriteJSONError(conn net.Conn, err error) error {
	return WriteJSON(conn, &StreamedResponse{
		MessageType: "error",
		MessageData: err.Error(),
		PercentDone: 100,
	})
}

func WriteJSON(conn net.Conn, data interface{}) error {
	w := wsutil.NewWriter(conn, ws.StateServerSide, ws.OpText)
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(&data); err != nil {
		return err
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}

func ReadJSON(conn net.Conn, data interface{}) error {
	r := wsutil.NewReader(conn, ws.StateServerSide)
	decoder := json.NewDecoder(r)

	hdr, err := r.NextFrame()
	if err != nil {
		return err
	}

	if hdr.OpCode == ws.OpClose {
		return io.EOF
	}

	if err := decoder.Decode(data); err != nil {
		return err
	}

	return nil
}

func Websockify(w http.ResponseWriter, req *http.Request) (net.Conn, error) {
	conn, _, _, err := ws.UpgradeHTTP(req, w)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
