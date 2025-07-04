import React, { useState } from 'react';
import { Card, Button, Badge, Typography, Row, Col, Divider } from 'antd';
import { CheckCircleOutlined, CrownOutlined, RocketOutlined, StarOutlined } from '@ant-design/icons';
import useBillingStore from '../store/billing.store';
import { usePlansQuery } from '../../plans/hooks/usePlanQuery';
import { useMailCalcQuery } from '../../../hooks/useMailDataQuery';
import { APIPlanData } from '../../plans/interface/plan.interface';

const { Title, Text, Paragraph } = Typography;

const PricingPlans: React.FC = () => {
    const { setPaymentValues, initializePayment } = useBillingStore();
    const [loadingPlanId, setLoadingPlanId] = useState<string | null>(null);
    const { data: planData } = usePlansQuery();
    const { data: mailData } = useMailCalcQuery();

    const currentPlan = mailData?.payload.plan;

    const handlePlanSelection = async (plan: APIPlanData) => {
        setPaymentValues({
            plan_id: plan.id,
            duration: plan.billing_cycle,
            amount_to_pay: plan.price,
            payment_method: "Paystack"
        });

        localStorage.setItem('planSubscription', JSON.stringify(plan));
        setLoadingPlanId(plan.id);

        try {
            await initializePayment();
        } finally {
            setLoadingPlanId(null);
        }
    };

    const getPlanIcon = (planName: string) => {
        const name = planName.toLowerCase();
        if (name.includes('enterprise')) return <CrownOutlined className="text-purple-600" />;
        if (name.includes('professional')) return <RocketOutlined className="text-blue-600" />;
        if (name.includes('basic')) return <StarOutlined className="text-green-600" />;
        return <CheckCircleOutlined className="text-gray-600" />;
    };

    const getPlanColor = (planName: string) => {
        const name = planName.toLowerCase();
        if (name.includes('enterprise')) return 'border-purple-200 hover:border-purple-400';
        if (name.includes('professional')) return 'border-blue-200 hover:border-blue-400';
        if (name.includes('basic')) return 'border-green-200 hover:border-green-400';
        return 'border-gray-200 hover:border-gray-400';
    };

    const getButtonStyle = (plan: APIPlanData) => {
        const isCurrentPlan = currentPlan === plan.name;
        // const isLoading = loadingPlanId === plan.id;
        const isFreePlan = plan.price === 0;

        if (isCurrentPlan) {
            return { type: 'default' as const, className: 'bg-gray-100 border-gray-300' };
        }
        if (isFreePlan) {
            return { type: 'default' as const, className: 'bg-green-50 border-green-200 text-green-700' };
        }

        const name = plan.name.toLowerCase();
        if (name.includes('enterprise')) {
            return { type: 'primary' as const, className: 'bg-gradient-to-r from-purple-600 to-purple-700 border-purple-600' };
        }
        if (name.includes('professional')) {
            return { type: 'primary' as const, className: 'bg-gradient-to-r from-blue-600 to-blue-700 border-blue-600' };
        }
        return { type: 'primary' as const, className: 'bg-gradient-to-r from-green-600 to-green-700 border-green-600' };
    };

    const renderPlanCard = (plan: APIPlanData) => {
        const isCurrentPlan = currentPlan === plan.name;
        const isLoading = loadingPlanId === plan.id;
        const isFreePlan = plan.price === 0;
        const buttonStyle = getButtonStyle(plan);

        return (
            <Col xs={24} sm={12} lg={6} key={plan.id}>
                <Card
                    className={`h-full transition-all duration-300 hover:shadow-xl ${getPlanColor(plan.name)} ${isCurrentPlan ? 'border-2 border-blue-500 shadow-lg' : ''
                        }`}
                    hoverable
                    style={{ borderRadius: '16px' }}
                >
                    <div className="text-center mb-6">
                        <div className="flex justify-center items-center mb-3">
                            {getPlanIcon(plan.name)}
                            <Title level={3} className="ml-2 mb-0">
                                {plan.name}
                            </Title>
                        </div>

                        {isCurrentPlan && (
                            <Badge
                                count="Current Plan"
                                className="mb-3"
                                style={{ backgroundColor: '#52c41a' }}
                            />
                        )}

                        <div className="mb-4">
                            <span className="text-4xl font-bold text-gray-800">
                                {plan.price === 0 ? 'Free' : `₦${plan.price.toLocaleString()}`}
                            </span>
                            <Text className="text-gray-500 ml-2">/{plan.billing_cycle}</Text>
                        </div>

                        <Paragraph className="text-gray-600 mb-6 min-h-[48px]">
                            {plan.description}
                        </Paragraph>
                    </div>

                    <Button
                        {...buttonStyle}
                        size="large"
                        block
                        loading={isLoading}
                        disabled={isCurrentPlan || (isFreePlan && !isCurrentPlan)}
                        onClick={() => handlePlanSelection(plan)}
                        className={`mb-6 font-semibold ${buttonStyle.className}`}
                        style={{ borderRadius: '8px', height: '48px' }}
                    >
                        {isLoading ? 'Processing...' :
                            isCurrentPlan ? 'Current Plan' :
                                isFreePlan ? 'Free Plan' :
                                    'Choose Plan'}
                    </Button>

                    <Divider />

                    <div className="space-y-3">
                        <Text strong className="text-gray-700">Features included:</Text>
                        {plan.features.map((feature, index) => (
                            <div key={index} className="flex items-start space-x-3">
                                <CheckCircleOutlined className="text-green-500 mt-1 flex-shrink-0" />
                                <div>
                                    <Text strong className="text-gray-800">{feature.name}:</Text>
                                    <Text className="text-gray-600 ml-2">{feature.value}</Text>
                                </div>
                            </div>
                        ))}

                        {/* Mailing Limits */}
                        <Divider className="my-4" />
                        <Text strong className="text-gray-700">Mailing Limits:</Text>
                        <div className="space-y-2">
                            <div className="flex items-center space-x-3">
                                <CheckCircleOutlined className="text-blue-500" />
                                <Text className="text-gray-600">
                                    <Text strong>Daily:</Text> {plan.mailing_limits.daily_limit === 0 ? 'Unlimited' : plan.mailing_limits.daily_limit.toLocaleString()}
                                </Text>
                            </div>
                            <div className="flex items-center space-x-3">
                                <CheckCircleOutlined className="text-blue-500" />
                                <Text className="text-gray-600">
                                    <Text strong>Monthly:</Text> {plan.mailing_limits.monthly_limit === 0 ? 'Unlimited' : plan.mailing_limits.monthly_limit.toLocaleString()}
                                </Text>
                            </div>
                            <div className="flex items-center space-x-3">
                                <CheckCircleOutlined className="text-blue-500" />
                                <Text className="text-gray-600">
                                    <Text strong>Max per email:</Text> {plan.mailing_limits.max_recipients_per_mail === 0 ? 'Unlimited' : plan.mailing_limits.max_recipients_per_mail.toLocaleString()}
                                </Text>
                            </div>
                        </div>
                    </div>
                </Card>
            </Col>
        );
    };

    return (
        <div className="min-h-screen bg-gradient-to-b from-gray-50 to-white py-12">
            <div className="container mx-auto px-4 max-w-7xl">
                <div className="text-center mb-12">
                    <Title level={1} className="text-5xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent mb-4">
                        Upgrade Your Marketing Plan
                    </Title>
                    <Text className="text-xl text-gray-600 max-w-2xl mx-auto">
                        Choose the perfect plan to supercharge your email marketing campaigns and grow your business
                    </Text>
                </div>

                <Row gutter={[24, 24]} className="justify-center">
                    {planData?.payload.map((plan: any) => renderPlanCard(plan))}
                </Row>

                <div className="text-center mt-12">
                    <Text className="text-gray-500">
                        Need a custom solution? <Text strong className="text-blue-600 cursor-pointer">Contact our sales team</Text>
                    </Text>
                </div>
            </div>
        </div>
    );
};

export default PricingPlans;
