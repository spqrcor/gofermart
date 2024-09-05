-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.withdraw_list (
                               id uuid DEFAULT gen_random_uuid() NOT NULL,
                               order_id uuid NOT NULL,
                               sum integer NOT NULL CHECK (sum > 0),
                               created_at timestamptz DEFAULT now() NOT NULL,
                               CONSTRAINT withdraw_list_pkey PRIMARY KEY (id)
);
CREATE INDEX withdraw_list_order_id_idx ON public.withdraw_list (order_id);
CREATE INDEX withdraw_list_created_at_idx ON public.withdraw_list (created_at DESC);

COMMENT ON COLUMN public.withdraw_list.id IS 'UUID';
COMMENT ON COLUMN public.withdraw_list.order_id IS 'UUID заказа';
COMMENT ON COLUMN public.withdraw_list.sum IS 'Сумма баллов к списанию';
COMMENT ON COLUMN public.withdraw_list.created_at IS 'Дата создания';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.withdraw_list;
-- +goose StatementEnd
