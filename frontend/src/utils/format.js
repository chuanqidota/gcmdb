/**
 * Format data as pretty-printed JSON string.
 * Handles both objects and JSON strings.
 */
export function formatJson(data, indent = 2) {
  try {
    if (typeof data === 'string') {
      return JSON.stringify(JSON.parse(data), null, indent)
    }
    return JSON.stringify(data, null, indent)
  } catch {
    return String(data)
  }
}
