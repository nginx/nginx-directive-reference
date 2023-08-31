import Directive from './Directive'
import Tree from './Tree'
import { useSearchParams } from './params'

function App() {
  const [{ modId, directive }] = useSearchParams()

  return (
    <>
      <h1>Reference Browser</h1>
      <div className="container">
        <div style={{ minWidth: '20rem' }}>
          <Tree modId={modId} />
        </div>
        <div style={{ flex: '1 1 auto' }}>
          {modId && directive && (
            <Directive modId={modId} directive={directive} />
          )}
        </div>
      </div>
    </>
  )
}

export default App
