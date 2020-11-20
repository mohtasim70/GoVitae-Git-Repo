import React from "react";
import "./App.css";
// import the Container Component from the semantic-ui-react
import { Container } from "semantic-ui-react";
// import the ToDoList component
import Blockchain from "./Blockchain";
function App() {
  return (
    <div>
      <Container>
        <Blockchain />
      </Container>
    </div>
  );
}
export default App;
