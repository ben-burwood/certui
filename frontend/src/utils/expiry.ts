import type { ExpiryCountdown } from "@/types/expiry";

// Utility function to calculate days until a given date
export function daysUntil(date: string, now: Date): number {
  const expiry = new Date(date);

  // Calculate difference in milliseconds
  const diffMs = expiry.getTime() - now.getTime();

  // Convert ms to days
  return Math.ceil(diffMs / (1000 * 60 * 60 * 24));
}

// Utility function to check if a given date is in the past
export function isExpired(date: string, now: Date): boolean {
  const expiry = new Date(date);
  return expiry.getTime() < now.getTime();
}

// Utility function to calculate days, hours, and minutes until a given date
export function expiryUntilCountdown(date: string, now: Date): ExpiryCountdown {
  const expiry = new Date(date);

  const diffMs = expiry.getTime() - now.getTime();
  const diffMinutes = Math.floor(diffMs / (1000 * 60));
  const days = Math.floor(diffMinutes / (60 * 24));
  const hours = Math.floor((diffMinutes % (60 * 24)) / 60);
  const minutes = diffMinutes % 60;
  return { days, hours, minutes };
}
