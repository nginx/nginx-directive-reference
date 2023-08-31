import { useEffect, useState } from 'react'

type SearchBoxProps = {
  onChange: (term: string) => void
  /** how many milliseconds to wait for user input to stop changing end before
   * calling onChange, defaults to 500ms */
  debounceMs?: number
}

export function SearchBox({ onChange, debounceMs = 500 }: SearchBoxProps) {
  const [term, setTerm] = useState<string>()
  useEffect(() => {
    if (term === undefined) {
      return
    }
    const tId = setTimeout(() => onChange(term.trim()), debounceMs)
    return () => clearTimeout(tId)
  }, [term, onChange, debounceMs])

  return (
    <div className="container">
      🔍
      <input
        type="text"
        placeholder="search directives"
        role="searchbox"
        value={term}
        onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
          setTerm(event.target.value)
        }}
      />
    </div>
  )
}
