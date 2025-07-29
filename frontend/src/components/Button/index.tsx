import { Button } from 'antd';
import React from 'react';
import { useNavigate } from 'react-router-dom';
import styled from 'styled-components';

interface BottonProps {
  className?: string;
  buttons?: Array<{
    text: string;
    onClick?: () => void;
    key: string;
    to?: string;
    className?: string;
    style?: React.CSSProperties;
    render?: () => React.ReactNode;
  }>;
}

const BottonWrapper = styled.div`
  display: flex;
  justify-content: center;
  gap: 20px;
  padding: 10px;
  border-top: 1px solid #1890ff;
`;

const Botton: React.FC<BottonProps> = ({ className, buttons = [] }) => {
  const navigate = useNavigate();

  return (
    <BottonWrapper className={className}>
      {buttons.map((button) =>
        button.render ? (
          <React.Fragment key={button.key}>{button.render()}</React.Fragment>
        ) : (
          <Button
            key={button.key}
            type="primary"
            className={button.className}
            style={button.style}
            onClick={() => {
              if (button.onClick) {
                button.onClick();
              } else if (button.to) {
                navigate(button.to);
              }
            }}
          >
            {button.text}
          </Button>
        ),
      )}
    </BottonWrapper>
  );
};

export default Botton;
