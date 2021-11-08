import React, { useState } from "react";
import { SubmitHandler } from "react-hook-form";
import styled from "styled-components";
import Form from "../components/form";
import axios from "axios";
import { ToastContainer, toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.min.css";

const Page = styled.main`
  color: #232129;
  padding: 96px;
  font-family: -apple-system, Roboto, sans-serif, serif;
  display: flex;
`;
const Heading = styled.h1`
  margin-top: 0;
  margin-bottom: ${({ short }) => (short ? "64px" : "32px")};
  max-width: ${({ short }) => (short ? "320px" : "auto")};
`;

const Accent = styled.span`
  color: #663399;
`;

const Recipe = styled.div`
  margin-left: 1.5rem;
`;
type Inputs = {
  pageUrl: string;
};

type Recipe = {
  ingredients?: string[];
  title?: string;
};
const IndexPage = () => {
  const [recipe, setRecipe] = useState<Recipe>({});

  const onFormSubmit: SubmitHandler<Inputs> = async (formData) => {
    try {
      const { data, status, ...response } = await axios.post(
        "/.netlify/functions/recipe",
        formData
      );
      if (status !== 200) {
        throw new Error("Unable to get recipe");
      }

      setRecipe(data);
      console.log(status, data, response);
    } catch (error) {
      toast.error(error.response.data);
      console.log(error);
    }
  };
  return (
    <Page>
      <title>Home Page</title>
      <div>
        <Heading short>
          Recipe Lens
          <br />
          <Accent>— cut the bullshit </Accent>
          <br />
          <span role="img" aria-label="Party popper emojis">
            🌮🥗🍛
          </span>
        </Heading>
        <Form onSubmit={onFormSubmit} />
      </div>
      <Recipe>
        {recipe.title ? <Heading>{recipe.title}</Heading> : null}
        {recipe?.ingredients?.length > 0 ? (
          <ul>
            {recipe?.ingredients?.map((ingridient) => (
              <li>{ingridient}</li>
            ))}
          </ul>
        ) : null}
        <ul></ul>
      </Recipe>
      <ToastContainer
        position="bottom-right"
        hideProgressBar
        autoClose={5000}
        newestOnTop
      />
    </Page>
  );
};

export default IndexPage;
