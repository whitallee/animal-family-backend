ALTER TABLE "animals" 
ADD COLUMN "isMemorialized" BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN "lastMessage" TEXT,
ADD COLUMN "memorialPhotos" JSONB;

