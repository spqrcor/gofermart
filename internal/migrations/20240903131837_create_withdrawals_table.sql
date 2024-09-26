-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.withdrawals (
                               id uuid DEFAULT gen_random_uuid() NOT NULL,
                               user_id uuid NOT NULL,
                               number varchar(20) NOT NULL,
                               sum numeric(12, 2) NOT NULL CHECK (sum > 0),
                               created_at timestamptz DEFAULT now() NOT NULL,
                               CONSTRAINT withdraw_list_pkey PRIMARY KEY (id),
                               CONSTRAINT withdrawals_user_id_fkey
                                   FOREIGN KEY (user_id)
                                   REFERENCES public.users(id)
                                   ON DELETE CASCADE
);
CREATE INDEX withdrawals_user_id_idx ON public.withdrawals (user_id);
CREATE INDEX withdrawals_created_at_idx ON public.withdrawals (created_at DESC);
CREATE UNIQUE INDEX withdrawals_number_idx ON public.withdrawals (number);

COMMENT ON COLUMN public.withdrawals.id IS 'UUID';
COMMENT ON COLUMN public.withdrawals.user_id IS 'UUID пользователя';
COMMENT ON COLUMN public.withdrawals.number IS 'Номер заказа';
COMMENT ON COLUMN public.withdrawals.sum IS 'Сумма баллов к списанию';
COMMENT ON COLUMN public.withdrawals.created_at IS 'Дата создания';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.withdrawals;
-- +goose StatementEnd
