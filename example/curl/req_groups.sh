#!/bin/sh

TOKEN=$(head -1 ./gitlab.token.txt| tr -d '\n')

group_id=11827547
#group_id=11536035

#curl --request GET --header "PRIVATE-TOKEN: $TOKEN" --header "Content-Type: application/json" \
#		"https://gitlab.com/api/v4/groups/$group_id/boards"\
#;
#
#board_id=2576399

curl --request GET --header "PRIVATE-TOKEN: $TOKEN" --header "Content-Type: application/json" \
		"https://gitlab.com/api/v4/groups/$group_id/issues?assignee_username=mxpatlas"\
;
