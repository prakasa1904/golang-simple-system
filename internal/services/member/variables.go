package member

var (
	SelectColumn         = "id, fullname, username, email, phone"
	SelectColumnWithJoin = "`member`.`id`,`member`.`fullname`,`member`.`username`,`member`.`phone`, `member`.`email`,`member`.`password`,`member`.`created_at`,`member`.`updated_at`,`member`.`group_id`"
	AllowedFilterQuery   = []string{"id", "fullname", "username"}
)
