import { useSearchParams as baseUseSearchParams } from 'react-router-dom'

interface SearchParams {
  modId?: string
  directive?: string
}

type SetSearchParams = (x: Required<SearchParams>) => void

/** type-safe wrapper around react-router-dom's useSearchParams for the search params we care about */
export function useSearchParams(): [SearchParams, SetSearchParams] {
  const [sp, setSP] = baseUseSearchParams()
  const modId = sp.get('modId') || undefined
  const directive = sp.get('directive') || undefined
  return [{ modId, directive }, (x: Required<SearchParams>) => setSP(x)]
}

export function toHref(sp: Required<SearchParams>): string {
  return '?' + new URLSearchParams(sp).toString()
}
