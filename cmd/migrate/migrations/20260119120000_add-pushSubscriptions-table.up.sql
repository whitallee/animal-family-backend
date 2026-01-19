CREATE TABLE IF NOT EXISTS "pushSubscriptions" (
    "subscriptionId" SERIAL PRIMARY KEY,
    "userId" INTEGER NOT NULL,
    "endpoint" TEXT NOT NULL,
    "p256dh" TEXT NOT NULL,
    "auth" TEXT NOT NULL,
    "userAgent" TEXT,
    "createdAt" TIMESTAMP DEFAULT NOW(),
    "lastUsed" TIMESTAMP DEFAULT NOW(),

    FOREIGN KEY ("userId") REFERENCES users("userId") ON DELETE CASCADE,
    UNIQUE("userId", "endpoint")
);

CREATE INDEX idx_push_subscriptions_user ON "pushSubscriptions"("userId");
