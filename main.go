package main

import (
	"io"

	"github.com/jackpal/bencode-go"
)

type bencodeInfo struct {
	Length      int
	Name        string
	Pieces      string
	PieceLength int
}

type bencodeTorrent struct {
	Announce string
	Info     bencodeInfo
}

func Open(r io.Reader) (*bencodeTorrent, error) {
	btorr := bencodeTorrent{}
	err := bencode.Unmarshal(r, &btorr)
	if err != nil {
		return nil, err
	}
	return &btorr, nil

}
