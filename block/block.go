package block

type Block struct {
	Data         string
	PreviousHash []byte
}

var isInitialized bool = false

func CreateRootBlock() *Block {
	if !isInitialized {
		isInitialized = true
		return &Block{Data: "Root", PreviousHash: nil}
	}

	return nil
}
