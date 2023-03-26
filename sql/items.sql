CREATE TABLE public."Items" (
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