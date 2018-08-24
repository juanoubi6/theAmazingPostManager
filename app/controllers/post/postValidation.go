package post

import "errors"

func validateTitle(title string) (bool, string) {
	if title == "" {
		return false, "Title cannot be empty"
	}

	if len(title) > 60 {
		return false, "Title must be less than 60 characters"
	}

	return true, ""

}

func validateDescription(description string) (bool, string) {
	if description == "" {
		return false, "Description cannot be empty"
	}

	if len(description) > 1000 {
		return false, "Description must be less than 1000 characters"
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