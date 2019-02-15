import * as React from "react";
import styled from "styled-components";
import Widget from "../elements/Widget";
import BlockButton from "../elements/BlockButton";
import Input from "./Input";

const Inputs = styled.div`
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
`;
const Label = styled.div`
    color: grey;
    font-size: 24px;
    text-transform: uppercase;
    border-bottom: 1px solid grey;
    box-sizing: border-box;
    padding: 5px;
    padding-bottom: 15px;
    margin-bottom: 30px;
`;
const Links = styled.div`
    display: flex;
    justify-content: space-between;
    text-decoration: none;
    color: rgba(30, 100, 200, 1);
    margin-bottom: 20px;
`;
const Link = styled.a`
    text-decoration: none;
`;
const Button = styled(BlockButton)`
    width: 100%;
`;

interface LinkParams {
    text: string;
    onClick: () => void;
}
interface ButtonParams {
    text: string;
    onClick: () => void;
}
interface AuthFormProps {
    title: string;
    emailValue: string;
    emailError: string | null;
    onEmailChange: (e: React.FormEvent<HTMLInputElement>) => void;
    passwordValue: string;
    passwordError: string | null;
    onPasswordChange: (e: React.FormEvent<HTMLInputElement>) => void;
    linkParams: LinkParams[];
    buttonParams: ButtonParams;
}
class AuthForm extends React.PureComponent<AuthFormProps, any> {
    render() {
        const props = this.props;
        return (
            <Widget>
                <Label>{props.title}</Label>
                <Inputs>
                    <Input
                        label={"email"}
                        type={"email"}
                        value={props.emailValue}
                        onChange={props.onEmailChange}
                        error={props.emailError}
                    />
                    <Input
                        label={"password"}
                        type={"password"}
                        value={props.passwordValue}
                        onChange={props.onPasswordChange}
                        error={props.passwordError}
                    />
                </Inputs>
                <Links>
                    {props.linkParams.map(param => (
                        <Link key={param.text} onClick={param.onClick}>
                            {param.text}
                        </Link>
                    ))}
                </Links>
                <Button onClick={props.buttonParams.onClick}>
                    {props.buttonParams.text}
                </Button>
            </Widget>
        );
    }
}
export default AuthForm;
