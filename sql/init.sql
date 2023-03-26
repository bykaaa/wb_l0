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



	CREATE TABLE "Delivery" (
	order_uid text NOT NULL,
	"name" text NOT NULL,
	phone text NOT NULL,
	zip text NOT NULL,
	city text NOT NULL,
	address text NOT NULL,
	region text NOT NULL,
	email text NOT NULL,
	CONSTRAINT "Delivery_order_uid_fkey" FOREIGN KEY (order_uid) REFERENCES public."Order"(order_uid) ON DELETE CASCADE
);




CREATE TABLE "Items" (
	order_uid text not null,
	id serial not null,
	chrt_id integer not null ,
	track_number text not null,
	price integer not null,
	rid text not null,
	name text not null,
	sale integer not null,
	size text not null,
	total_price integer not null,
	nm_id integer not null,
	brand text not null,
	status integer not null,
	FOREIGN KEY (order_uid) REFERENCES public."Order" (order_uid) ON DELETE CASCADE
	);


	CREATE TABLE "Payment" (
	"transaction" text NOT NULL,
	request_id text NOT NULL,
	currency varchar(25) NOT NULL,
	provider text NOT NULL,
	amount int4 NOT NULL,
	payment_dt int4 NOT NULL,
	bank text NOT NULL,
	delivery_cost int4 NOT NULL,
	goods_total int4 NOT NULL,
	custom_fee int4 NOT NULL,
	CONSTRAINT "Payment_transaction_fkey" FOREIGN KEY ("transaction") REFERENCES public."Order"(order_uid) ON DELETE CASCADE
);

