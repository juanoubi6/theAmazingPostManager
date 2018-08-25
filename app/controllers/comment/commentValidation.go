package comment

import (
	"errors"
)

func validateMessage(message string) (bool, string) {
	if message == "" {
		return false, "Message cannot be empty"
	}

	if len(message) > 160 {
		return false, "Message must be less than 160 characters"
	}

	return true, ""

}

func getVoteValue(vote string)(bool,error){

	if vote == ""{
		return false, errors.New("Vote is empty")
	}

	if vote == "true"{
		return true,nil
	}else if vote == "false"{
		return false,nil
	}else{
		return false,errors.New("Invalid vote")
	}

}