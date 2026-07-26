/// App-wide constants and configuration values.
///
/// In production, load these from environment variables or a secure vault.
class AppConstants {
  AppConstants._();

  // ── Supabase ──────────────────────────────────────────────
  // Replace with your Supabase project URL and anon key.
  static const String supabaseUrl = String.fromEnvironment(
    'SUPABASE_URL',
    defaultValue: 'https://YOUR_PROJECT.supabase.co',
  );

  static const String supabaseAnonKey = String.fromEnvironment(
    'SUPABASE_ANON_KEY',
    defaultValue: 'YOUR_ANON_KEY',
  );

  // ── App ───────────────────────────────────────────────────
  static const String appName = '{{project_name}}';
  static const String appVersion = '1.0.0';

  // ── UPI ───────────────────────────────────────────────────
  // UPI ID for receiving payments. Replace with your VPA.
  static const String defaultUpiId = 'your-vpa@upi';

  // ── Limits ────────────────────────────────────────────────
  static const int maxGroupMembers = 50;
  static const int maxExpenseAmount = 1000000; // in INR
}
