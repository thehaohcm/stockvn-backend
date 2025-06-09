CREATE DATABASE stockvn_db;

USING stockvn_db;

CREATE TABLE public.symbols_watchlist (
	symbol varchar NOT NULL,
	highest_price int8 NULL,
	lowest_price int8 NULL,
	auto_trade bool NOT NULL DEFAULT false,
	updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT symbols_watchlist_pkey PRIMARY KEY (symbol)
);

INSERT INTO public.symbols_watchlist (symbol,highest_price,lowest_price,auto_trade,updated_at) VALUES
	 ('VIC',98400,39900,false,'2025-06-09 13:03:11.399154+07'),
	 ('HHS',15150,6940,false,'2025-06-09 13:03:11.399154+07'),
	 ('CRC',10600,6260,false,'2025-06-09 13:03:11.399154+07'),
	 ('DXS',8850,5210,false,'2025-06-09 13:03:11.399154+07'),
	 ('SHB',13392,8767,false,'2025-06-09 13:03:11.399154+07'),
	 ('DC4',14200,8697,false,'2025-06-09 13:03:11.399154+07'),
	 ('GEG',16850,10650,false,'2025-06-09 13:03:11.399154+07'),
	 ('BAF',36600,17300,false,'2025-06-09 13:03:11.399154+07'),
	 ('VHM',77600,34500,false,'2025-06-09 13:03:11.399154+07'),
	 ('REE',78000,60663,false,'2025-06-09 13:03:11.399154+07');
