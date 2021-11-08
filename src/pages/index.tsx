import * as React from "react";
import styled from "styled-components";
import Form from "../components/form";

const Page = styled.main`
  color: #232129;
  padding: 96px;
  font-family: -apple-system, Roboto, sans-serif, serif;
`;
const Heading = styled.h1`
  margin-top: 0;
  margin-bottom: 64px;
  max-width: 320px;
`;

const Accent = styled.span`
  color: #663399;
`;
const IndexPage = () => {
  return (
    <Page>
      <title>Home Page</title>
      <Heading>
        Congratulations
        <br />
        <Accent>— you just made a Gatsby site! </Accent>
        <span role="img" aria-label="Party popper emojis">
          🎉🎉🎉
        </span>
      </Heading>
      <Form />
    </Page>
  );
};

export default IndexPage;
