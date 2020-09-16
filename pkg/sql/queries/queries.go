package queries

const ALL_ITEMS = `
	SELECT item_id, model_id, item_category_id, page_no, item_no, foreign_id, item_name, price FROM item`

const ITEM_DETAILS = `
	SELECT I.id, I.item_id, I.model_id, M.name AS model_name, I.item_category_id, IC.name AS item_category_name, I.page_no, I.item_no, I.foreign_id, I.item_name, I.price
	FROM item I 
	LEFT JOIN model M ON M.id = I.model_id
	LEFT JOIN item_category IC ON IC.id = I.item_category_id
	WHERE I.item_id = ?;
`

const CHECK_IF_NUMBER_EXISTS = `
	SELECT COUNT(U.id) AS exist
	FROM user U 
	WHERE U.country_code = ? AND U.number = ?
`

const CHECK_IF_PIN_CORRECT = `
	SELECT COUNT(id) AS valid
	FROM user 
	WHERE country_code = ? AND number = ? AND pin = ?
`

const CHECK_IF_HASH_CORRECT = `
	SELECT COUNT(id) AS valid
	FROM user 
	WHERE country_code = ? AND number = ? AND hash = ?
`

const USER_DETAILS = `
SELECT first_name, last_name, display_name, dob, status
FROM user WHERE number = ?
`
