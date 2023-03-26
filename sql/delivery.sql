CREATE TABLE public."Delivery" (
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