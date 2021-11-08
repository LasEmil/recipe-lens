import axios from "axios";
import React from "react";
import { useForm, SubmitHandler } from "react-hook-form";
import styled from "styled-components";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";

const Input = styled.input`
  padding: 0.75rem 1rem;
`;

const handleColorType = (color) => {
  switch (color) {
    case "primary":
      return "color: #663399; background: #eee";
    case "danger":
      return "color: #fff; background: #f56342;";
    default:
      return "color: #000; background: #eee;";
  }
};

const Button = styled.button`
  display: block;
  margin: 5px 0;
  cursor: pointer;
  border: 0;
  font-size: 1.25rem;
  padding: 0.375rem 1rem;
  font-weight: 600;
  ${({ color }) => handleColorType(color)};

  &:focus {
    outline: 0;
  }
`;

type Inputs = {
  pageUrl: string;
};

const Schema = yup
  .object({
    pageUrl: yup.string().url().required(),
  })
  .required();

const Form = () => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<Inputs>({ resolver: yupResolver(Schema) });
  const onSubmit: SubmitHandler<Inputs> = async (data) => {
    try {
      const response = await axios.post("/.netlify/functions/hello", data);
      console.log(response);
    } catch (error) {
      console.log(error);
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <Input {...register("pageUrl", { required: true })} />
      <p>{errors.pageUrl?.message}</p>

      <Button color="primary" as="input" type="submit" />
    </form>
  );
};

export default Form;
