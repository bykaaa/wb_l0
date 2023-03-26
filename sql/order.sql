CREATE TABLE "Order" (
	order_uid text primary key unique not null,
	track_number text not null,
	entry text not null,
	locale varchar(25) not null,
	internal_signature text not null,
	customer_id text not null,
	delivery_service text not null,
	shardkey text not null,
	sm_id integer not null,
	date_created timestamp not null,
	oof_shard text not null
	);