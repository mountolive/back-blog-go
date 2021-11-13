import { useDeno } from 'aleph/react';

export default function config(): { [env: string]: string } {
  const values = useDeno(() => Deno.env.toObject());
  return values;
}
