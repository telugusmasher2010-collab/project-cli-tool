-- {{project_name}} — Supabase Schema
-- Run this in the Supabase SQL Editor to create all tables.

-- ── Groups ──────────────────────────────────────────────────
create table if not exists groups (
  id         uuid primary key default gen_random_uuid(),
  name       text not null,
  created_by uuid references auth.users(id) not null,
  created_at timestamptz default now()
);

-- ── Group Members ───────────────────────────────────────────
create table if not exists group_members (
  id        uuid primary key default gen_random_uuid(),
  group_id  uuid references groups(id) on delete cascade not null,
  user_id   uuid references auth.users(id) not null,
  joined_at timestamptz default now(),
  unique(group_id, user_id)
);

-- ── Expenses ────────────────────────────────────────────────
create table if not exists expenses (
  id          uuid primary key default gen_random_uuid(),
  group_id    uuid references groups(id) on delete cascade not null,
  paid_by     uuid references auth.users(id) not null,
  title       text not null,
  amount      numeric(12,2) not null check (amount > 0),
  currency    text default 'INR',
  split_type  text default 'equal' check (split_type in ('equal', 'custom')),
  created_at  timestamptz default now()
);

-- ── Expense Splits ──────────────────────────────────────────
create table if not exists expense_splits (
  id         uuid primary key default gen_random_uuid(),
  expense_id uuid references expenses(id) on delete cascade not null,
  user_id    uuid references auth.users(id) not null,
  amount     numeric(12,2) not null check (amount >= 0),
  settled    boolean default false,
  settled_at timestamptz,
  unique(expense_id, user_id)
);

-- ── Indexes ─────────────────────────────────────────────────
create index if not exists idx_group_members_group on group_members(group_id);
create index if not exists idx_expenses_group on expenses(group_id);
create index if not exists idx_expense_splits_expense on expense_splits(expense_id);
create index if not exists idx_expense_splits_user on expense_splits(user_id);

-- ── Row Level Security ──────────────────────────────────────
alter table groups enable row level security;
alter table group_members enable row level security;
alter table expenses enable row level security;
alter table expense_splits enable row level security;

-- Users can read groups they belong to.
create policy "Members can view groups"
  on groups for select
  using (
    id in (
      select group_id from group_members
      where user_id = auth.uid()
    )
  );

-- Users can insert groups (becomes the creator).
create policy "Users can create groups"
  on groups for insert
  with check (created_by = auth.uid());

-- Members can view expenses in their groups.
create policy "Members can view expenses"
  on expenses for select
  using (
    group_id in (
      select group_id from group_members
      where user_id = auth.uid()
    )
  );

-- Members can insert expenses in their groups.
create policy "Members can insert expenses"
  on expenses for insert
  with check (
    group_id in (
      select group_id from group_members
      where user_id = auth.uid()
    )
  );

-- Members can view splits for expenses in their groups.
create policy "Members can view splits"
  on expense_splits for select
  using (
    expense_id in (
      select e.id from expenses e
      join group_members gm on gm.group_id = e.group_id
      where gm.user_id = auth.uid()
    )
  );
