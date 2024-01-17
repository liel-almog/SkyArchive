import { GlobalProvider } from "./context/GlobalProvider";
import { Router } from "./router/Router";

function App() {
  return (
    <GlobalProvider>
      <Router />
    </GlobalProvider>
  );
}

export default App;
