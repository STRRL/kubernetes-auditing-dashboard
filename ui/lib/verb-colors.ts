export const VERB_COLORS: Record<string, string> = {
  create: 'bg-green-500',
  update: 'bg-indigo-500',
  patch: 'bg-blue-500',
  delete: 'bg-red-500',
  deletecollection: 'bg-red-600',
  get: 'bg-gray-400',
  list: 'bg-gray-500',
  watch: 'bg-yellow-500',
};

export function getVerbColor(verb: string): string {
  return VERB_COLORS[verb.toLowerCase()] || 'bg-purple-500';
}

export function getVerbLabel(verb: string, uppercase: boolean = false): string {
  if (uppercase) {
    return verb.toUpperCase();
  }
  return verb;
}
