-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.users (
                                 id uuid DEFAULT gen_random_uuid() NOT NULL,
                                 login varchar(255) NOT NULL,
                                 password varchar(60) NOT NULL,
                                 balance integer DEFAULT 0 NOT NULL CHECK (balance >= 0),
                                 created_at timestamptz DEFAULT now() NOT NULL,
                                 CONSTRAINT users_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX users_login_idx ON public.users (login);

COMMENT ON COLUMN public.users.id IS 'UUID';
COMMENT ON COLUMN public.users.login IS 'Логин';
COMMENT ON COLUMN public.users.password IS 'Пароль';
COMMENT ON COLUMN public.users.balance IS 'Текущий баланс';
COMMENT ON COLUMN public.users.created_at IS 'Дата создания';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.users;
-- +goose StatementEnd
