package snowFlake

import "github.com/bwmarrin/snowflake"

func GenSeqId(nodeIndex int64) (*snowflake.ID, error) {

	node, err := snowflake.NewNode(nodeIndex)
	if err != nil {
		return nil, err
	}
	// Generate a snowflake ID.
	id := node.Generate()
	return &id, nil
}
