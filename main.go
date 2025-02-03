package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/jackpal/bencode-go"
)

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func (i *bencodeInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)

	if err != nil {
		return [20]byte{}, err
	}

	return sha1.Sum(buf.Bytes()), nil
}

func (i *bencodeInfo) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20
	buf := []byte(i.Pieces)

	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("recieved malformed pieces of length %d", len(buf))
		return nil, err
	}

	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)
	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}

	return hashes, nil
}

// Converts bencode structs into a TorrentFile struct
func (bto bencodeTorrent) toTorrentFile() (TorrentFile, error) {
	infoHash, err := bto.Info.hash()
	if err != nil {
		return TorrentFile{}, err
	}
	piecesHashes, err := bto.Info.splitPieceHashes()

	if err != nil {
		return TorrentFile{}, err
	}

	t := TorrentFile{
		Announce:    bto.Announce,
		InfoHash:    infoHash,
		PieceHashes: piecesHashes,
		PieceLength: bto.Info.PieceLength,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
	}

	return t, nil
}

func main() {
	file, err := os.Open("debian-12.9.0-amd64-netinst.iso.torrent")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	benInfo := bencodeTorrent{}

	err = bencode.Unmarshal(file, &benInfo)

	if err != nil {
		panic(err)
	}

	tFile, _ := benInfo.toTorrentFile()

	fmt.Println(tFile.Announce)

}
