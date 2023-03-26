CREATE TABLE public."Payment" (
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

