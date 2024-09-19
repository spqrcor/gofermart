-- +goose Up
-- +goose StatementBegin
CREATE TYPE e_order_status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE public.orders (
                              id uuid DEFAULT gen_random_uuid() NOT NULL,
                              user_id uuid NOT NULL,
                              number varchar(20) NOT NULL,
                              status e_order_status DEFAULT 'NEW' NOT NULL,
                              accrual numeric(12, 2),
                              created_at timestamptz DEFAULT now() NOT NULL,
                              updated_at timestamptz,
                              CONSTRAINT orders_pkey PRIMARY KEY (id),
                              CONSTRAINT orders_user_id_fkey
                                  FOREIGN KEY (user_id)
                                  REFERENCES public.users(id)
                                  ON DELETE CASCADE
);
CREATE UNIQUE INDEX orders_number_idx ON public.orders (number);
CREATE INDEX orders_user_id_idx ON public.orders (user_id);
CREATE INDEX orders_created_at_idx ON public.orders (created_at DESC);

COMMENT ON COLUMN public.orders.id IS 'UUID';
COMMENT ON COLUMN public.orders.user_id IS 'UUID пользователя';
COMMENT ON COLUMN public.orders.number IS 'Номер заказа';
COMMENT ON COLUMN public.orders.status IS 'Статус заказа';
COMMENT ON COLUMN public.orders.accrual IS 'Рассчитанные баллы к начислению';
COMMENT ON COLUMN public.orders.created_at IS 'Дата создания';
COMMENT ON COLUMN public.orders.updated_at IS 'Дата последнего обновления';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.orders;

DROP TYPE public.e_order_status;
-- +goose StatementEnd
