import React, { Component } from "react";
import axios from "axios";
import { Card, Header, Form, Input, Icon,Button } from "semantic-ui-react";

let endpoint = "http://localhost:8080";

class Blockchain extends Component {
  constructor(props) {
    super(props);

    this.state = {
      cCode: "",
      cName: "",
      cHrs: "",
      cGrade: "",
      items: []
    };
  }

  componentDidMount() {
    this.getTask()
  }

  onChange = event => {
    this.setState({
      [event.target.name]: event.target.value,
    });
  };

  onSubmit = () => {
    //console.log("pRINTING task", this.state.cHrs);
    if (this.state.cCode) {
      axios
        .post(
          endpoint + "/api/block",
          {
            cCode: this.state.cCode,
            cName: this.state.cName,
            cHrs: this.state.cHrs,
            cGrade: this.state.cGrade,
          },
          {
            headers: {
              "Content-Type": "application/x-www-form-urlencoded"
            }
          }
        )
        .then(res => {
          this.getTask();
          this.setState({
            cCode: "",
            cName: "",
            cHrs: "",
            cGrade: ""
          });
          console.log(res);
        });
    }
  };

  getTask = () => {
  axios.get(endpoint + "/api/block").then(res => {
    console.log(res);
    if (res.data) {
      this.setState({
        items: res.data.map(item => {
          return (
            <Card key={item.currHash} color="green" fluid style={{marginTop: "50px"}}>
              <Card.Content>
                <Card.Header textAlign="left">
                  <div style={{ wordWrap: "break-word", textAlign: "center"}}>Block Number: {item.blockno}</div>
                </Card.Header>

                <Card.Meta textAlign="left">
                  <div style={{ wordWrap: "break-word", marginTop: "20px" }}>Current Hash: {item.currHash}</div>
                  <div style={{ wordWrap: "break-word" }}>Previous Hash: {item.prevHash}</div>
                  <div style={{ wordWrap: "break-word", fontWeight: "bold", color: "black", fontSize: "20px", marginBottom:"20px", marginTop: "20px"}}>Course Details:</div>
                  <span style={{ paddingRight: 10 }}>Course Code: {item.course.cCode}</span>
                  <span style={{ paddingRight: 10 }}>Course Name: {item.course.cName}</span>
                  <span style={{ paddingRight: 10 }}>Course Credit Hours: {item.course.cHrs}</span>
                  <span style={{ paddingRight: 10 }}>Course Grade: {item.course.cGrade}</span>
                </Card.Meta>
              </Card.Content>
            </Card>
          );
        })
      });
    } else {
      this.setState({
        items: []
      });
    }
  });
};
  render() {
    return (
      <div>
        <div className="row">
          <Header className="header" as="h2" style={{textAlign:"center", marginTop:"20px", fontSize:"30px"}}>
            GoVitae Website
          </Header>
          <Header className="header" as="h2" style={{marginBottom:"20px"}}>
            Enter Course Details:
          </Header>
        </div>
        <div className="row">
          <Form onSubmit={this.onSubmit}>
            <Input
              type="text"
              name="cCode"
              onChange={this.onChange}
              value={this.state.cCode}
              fluid
              placeholder="Enter Course Code"
            />
            <Input
              type="text"
              name="cName"
              onChange={this.onChange}
              value={this.state.cName}
              fluid
              placeholder="Enter Course Name"
            />
            <Input
              type="text"
              name="cHrs"
              onChange={this.onChange}
              value={this.state.cHrs}
              fluid
              placeholder="Enter Credit Hours"
            />
            <Input
              type="text"
              name="cGrade"
              onChange={this.onChange}
              value={this.state.cGrade}
              fluid
              placeholder="Enter Grade"
            />
            <Button style={{backgroundColor: "black", color:"white", marginTop:"20px"}}>Create Block</Button>
          </Form>
        </div>
        <div className="row">
          <Card.Group>{this.state.items}</Card.Group>
        </div>
      </div>
    );
  }
}

export default Blockchain;
