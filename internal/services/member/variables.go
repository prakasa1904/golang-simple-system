package member

var (
	SelectColumn         = "id, fullname, username, email"
	SelectColumnWithJoin = "`member`.`id`,`member`.`fullname`,`member`.`username`,`member`.`email`,`member`.`password`,`member`.`created_at`,`member`.`updated_at`,`member`.`group_id`"
	AllowedFilterQuery   = []string{"id", "fullname", "username"}
)
