ALTER table IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "owner_currency_key";
ALTER table IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";
DROP table IF EXISTS users ;