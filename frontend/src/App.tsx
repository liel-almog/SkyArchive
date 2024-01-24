import { GlobalProvider } from "./context/GlobalProvider";
import { Router } from "./router/Router";
import "../src/utils/zod/errorMap";

function App() {
  return (
    <GlobalProvider>
      <Router />
    </GlobalProvider>
  );
}

export default App;
