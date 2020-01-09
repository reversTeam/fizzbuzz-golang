package services

import (
	"context"
	pb "github.com/reversTeam/fizzbuzz-golang/src/client/protobuf"
	"strconv"
)

type Client struct{}

func NewServer() *Client {
	return &Client{}
}

func (o *Client) FizzBuzz(ctx context.Context, in *pb.FizzBuzzRequest) (*pb.FizzBuzzResponse, error) {
	results := []string{}
	limit := int(in.Limit)
	fizz := int(in.Int1)
	buzz := int(in.Int2)
	fizzbuzz := in.Str1+in.Str2
	for i := 1; i <= limit; i++ {
		if i%(fizz*buzz) == 0 {
			results = append(results, fizzbuzz)
		} else if i%fizz == 0 {
			results = append(results, in.Str1)
		} else if i%buzz == 0 {
			results = append(results, in.Str2)
		} else {
			results = append(results, strconv.Itoa(i))
		}
	}
	return &pb.FizzBuzzResponse{Items: results}, nil
}