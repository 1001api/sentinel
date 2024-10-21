-- name: CheckUserHasActiveSub :one
SELECT EXISTS(
	SELECT 1 FROM subscriptions 
	WHERE user_id = $1
	AND next_billing_date > NOW() AND cancelled_at IS NULL
);
